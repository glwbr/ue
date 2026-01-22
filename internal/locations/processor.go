package locations

import (
	"fmt"
	"math"
	"slices"
	"time"
)

type Processor struct {
	registry              *Registry
	addressToVariantCount map[string]int
}

func NewProcessor(registry *Registry) *Processor {
	if registry == nil {
		registry = &Registry{Locations: []Location{}, NextID: 1}
	}

	p := &Processor{
		registry:              registry,
		addressToVariantCount: make(map[string]int),
	}
	p.rebuildVariantCounts()
	return p
}

func (p *Processor) Registry() *Registry {
	return p.registry
}

func (p *Processor) FindOrCreateLocation(address string, lat, lon float64) string {
	normalizedAddr := NormalizeAddress(address)

	if lat == 0 && lon == 0 {
		return ""
	}

	now := time.Now()

	if i := p.findLocationByAddress(normalizedAddr); i != -1 {
		loc := &p.registry.Locations[i]
		p.updateLocation(loc, normalizedAddr, lat, lon, now)
		return loc.ID
	}

	for i := range p.registry.Locations {
		loc := &p.registry.Locations[i]

		distance := HaversineDistance(lat, lon, loc.AvgLat, loc.AvgLon)
		if distance <= meterThreshold {
			p.updateLocation(loc, normalizedAddr, lat, lon, now)
			return loc.ID
		}
	}

	return p.createNewLocation(normalizedAddr, lat, lon, now)
}

func (p *Processor) rebuildVariantCounts() {
	p.addressToVariantCount = make(map[string]int)
	for _, loc := range p.registry.Locations {
		for _, variant := range loc.AddressVariants {
			p.addressToVariantCount[variant]++
		}
		p.addressToVariantCount[loc.CanonicalAddress]++
	}
}

func (p *Processor) findLocationByAddress(normalizedAddress string) int {
	for i, loc := range p.registry.Locations {
		if loc.CanonicalAddress == normalizedAddress {
			return i
		}
		if slices.Contains(loc.AddressVariants, normalizedAddress) {
			return i
		}
	}
	return -1
}

func (p *Processor) updateLocation(loc *Location, address string, lat, lon float64, now time.Time) {
	p.decrementAddressCount(loc.CanonicalAddress)

	loc.VisitCount++
	loc.LastSeen = now

	p.addAddressVariant(loc, address)
	p.incrementAddressCount(loc.CanonicalAddress)

	p.updateAverageCoordinates(loc, lat, lon)
	p.updateCanonicalAddress(loc, address)
}

func (p *Processor) addAddressVariant(loc *Location, address string) {
	p.incrementAddressCount(address)

	if slices.Contains(loc.AddressVariants, address) {
		return
	}

	if address != loc.CanonicalAddress {
		loc.AddressVariants = append(loc.AddressVariants, address)
	}
}

func (p *Processor) updateAverageCoordinates(loc *Location, lat, lon float64) {
	if lat == 0 && lon == 0 {
		return
	}

	totalVisits := float64(loc.VisitCount)
	loc.AvgLat = (loc.AvgLat*(totalVisits-1) + lat) / totalVisits
	loc.AvgLon = (loc.AvgLon*(totalVisits-1) + lon) / totalVisits
}

func (p *Processor) updateCanonicalAddress(loc *Location, address string) {
	if p.addressToVariantCount[address] > p.addressToVariantCount[loc.CanonicalAddress] {
		loc.CanonicalAddress = address
	}
}

func (p *Processor) incrementAddressCount(address string) {
	p.addressToVariantCount[address]++
}

func (p *Processor) decrementAddressCount(address string) {
	p.addressToVariantCount[address]--
}

func (p *Processor) createNewLocation(address string, lat, lon float64, now time.Time) string {
	locID := fmt.Sprintf("loc-%d", p.registry.NextID)
	p.registry.NextID++

	newLoc := Location{
		ID:               locID,
		CanonicalAddress: address,
		AddressVariants:  []string{},
		AvgLat:           lat,
		AvgLon:           lon,
		VisitCount:       1,
		FirstSeen:        now,
		LastSeen:         now,
	}

	p.registry.Locations = append(p.registry.Locations, newLoc)
	p.addressToVariantCount[address]++

	return locID
}

func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	lat1Rad := toRadians(lat1)
	lat2Rad := toRadians(lat2)
	dlatRad := toRadians(lat2 - lat1)
	dlonRad := toRadians(lon2 - lon1)

	a := math.Sin(dlatRad/2)*math.Sin(dlatRad/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlonRad/2)*math.Sin(dlonRad/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c * 1000
}

const (
	earthRadiusKm = 6371.0
)

func toRadians(v float64) float64 {
	return v * math.Pi / 180
}
