package parking

type Vehicle interface {
	Size() VehicleSize
	Name() string
}

type Bike struct {
	size VehicleSize
	name string
}

func (b Bike) Size() VehicleSize {
	return b.size
}

func (b Bike) Name() string {
	return b.name
}

func NewBike(name string) Bike {
	return Bike{size: SMALL, name: name}
}

type Car struct {
	size VehicleSize
	name string
}

func (c Car) Size() VehicleSize {
	return c.size
}

func (c Car) Name() string {
	return c.name
}

func NewCar(name string) Car {
	return Car{size: MEDIUM, name: name}
}

type Truck struct {
	size VehicleSize
	name string
}

func (t Truck) Size() VehicleSize {
	return t.size
}

func (t Truck) Name() string {
	return t.name
}

func NewTruck(name string) Truck {
	return Truck{size: LARGE, name: name}
}
