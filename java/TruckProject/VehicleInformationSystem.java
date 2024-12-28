import java.util.Scanner;

// Interface for all vehicles
interface Vehicle {
    String getMake();
    String getModel();
    int getYearOfManufacture();
}

// Interface for car-specific details
interface CarVehicle {
    void setNumberOfDoors(int numberOfDoors);
    int getNumberOfDoors();
    void setFuelType(String fuelType);
    String getFuelType();
}

// Interface for motorcycle-specific details
interface MotorVehicle {
    void setNumberOfWheels(int numberOfWheels);
    int getNumberOfWheels();
    void setMotorcycleType(String motorcycleType);
    String getMotorcycleType();
}

// Interface for truck-specific details
interface TruckVehicle {
    void setCargoCapacity(double cargoCapacity);
    double getCargoCapacity();
    void setTransmissionType(String transmissionType);
    String getTransmissionType();
}

// Car class implementing Vehicle and CarVehicle interfaces
class Car implements Vehicle, CarVehicle {
    private String make;
    private String model;
    private int yearOfManufacture;
    private int numberOfDoors;
    private String fuelType;

    public Car(String make, String model, int yearOfManufacture) {
        this.make = make;
        this.model = model;
        this.yearOfManufacture = yearOfManufacture;
    }

    @Override
    public String getMake() {
        return make;
    }

    @Override
    public String getModel() {
        return model;
    }

    @Override
    public int getYearOfManufacture() {
        return yearOfManufacture;
    }

    @Override
    public void setNumberOfDoors(int numberOfDoors) {
        this.numberOfDoors = numberOfDoors;
    }

    @Override
    public int getNumberOfDoors() {
        return numberOfDoors;
    }

    @Override
    public void setFuelType(String fuelType) {
        this.fuelType = fuelType;
    }

    @Override
    public String getFuelType() {
        return fuelType;
    }
}

// Motorcycle class implementing Vehicle and MotorVehicle interfaces
class Motorcycle implements Vehicle, MotorVehicle {
    private String make;
    private String model;
    private int yearOfManufacture;
    private int numberOfWheels;
    private String motorcycleType;

    public Motorcycle(String make, String model, int yearOfManufacture) {
        this.make = make;
        this.model = model;
        this.yearOfManufacture = yearOfManufacture;
    }

    @Override
    public String getMake() {
        return make;
    }

    @Override
    public String getModel() {
        return model;
    }

    @Override
    public int getYearOfManufacture() {
        return yearOfManufacture;
    }

    @Override
    public void setNumberOfWheels(int numberOfWheels) {
        this.numberOfWheels = numberOfWheels;
    }

    @Override
    public int getNumberOfWheels() {
        return numberOfWheels;
    }

    @Override
    public void setMotorcycleType(String motorcycleType) {
        this.motorcycleType = motorcycleType;
    }

    @Override
    public String getMotorcycleType() {
        return motorcycleType;
    }
}

// Truck class implementing Vehicle and TruckVehicle interfaces
class Truck implements Vehicle, TruckVehicle {
    private String make;
    private String model;
    private int yearOfManufacture;
    private double cargoCapacity;
    private String transmissionType;

    public Truck(String make, String model, int yearOfManufacture) {
        this.make = make;
        this.model = model;
        this.yearOfManufacture = yearOfManufacture;
    }

    @Override
    public String getMake() {
        return make;
    }

    @Override
    public String getModel() {
        return model;
    }

    @Override
    public int getYearOfManufacture() {
        return yearOfManufacture;
    }

    @Override
    public void setCargoCapacity(double cargoCapacity) {
        this.cargoCapacity = cargoCapacity;
    }

    @Override
    public double getCargoCapacity() {
        return cargoCapacity;
    }

    @Override
    public void setTransmissionType(String transmissionType) {
        this.transmissionType = transmissionType;
    }

    @Override
    public String getTransmissionType() {
        return transmissionType;
    }
}

public class VehicleInformationSystem {
    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);

        // Create a Car
        System.out.println("Enter Car Details:");
        System.out.print("Make: ");
        String carMake = scanner.nextLine();
        System.out.print("Model: ");
        String carModel = scanner.nextLine();
        System.out.print("Year of Manufacture: ");
        int carYear = scanner.nextInt();
        scanner.nextLine(); // Consume newline
        System.out.print("Number of Doors: ");
        int carDoors = scanner.nextInt();
        scanner.nextLine(); // Consume newline
        System.out.print("Fuel Type (petrol/diesel/electric): ");
        String carFuel = scanner.nextLine();

        Car car = new Car(carMake, carModel, carYear);
        car.setNumberOfDoors(carDoors);
        car.setFuelType(carFuel);

        // Create a Motorcycle
        System.out.println("Enter Motorcycle Details:");
        System.out.print("Make: ");
        String motorcycleMake = scanner.nextLine();
        System.out.print("Model: ");
        String motorcycleModel = scanner.nextLine();
        System.out.print("Year of Manufacture: ");
        int motorcycleYear = scanner.nextInt();
        scanner.nextLine(); // Consume newline
        System.out.print("Number of Wheels: ");
        int motorcycleWheels = scanner.nextInt();
        scanner.nextLine(); // Consume newline
        System.out.print("Motorcycle Type (sport/cruiser/off-road): ");
        String motorcycleType = scanner.nextLine();

        Motorcycle motorcycle = new Motorcycle(motorcycleMake, motorcycleModel, motorcycleYear);
        motorcycle.setNumberOfWheels(motorcycleWheels);
        motorcycle.setMotorcycleType(motorcycleType);

        // Create a Truck
        System.out.println("Enter Truck Details:");
        System.out.print("Make: ");
        String truckMake = scanner.nextLine();
        System.out.print("Model: ");
        String truckModel = scanner.nextLine();
        System.out.print("Year of Manufacture: ");
        int truckYear = scanner.nextInt();
        scanner.nextLine(); // Consume newline
        System.out.print("Cargo Capacity (tons): ");
        double truckCargoCapacity = scanner.nextDouble();
        scanner.nextLine(); // Consume newline
        System.out.print("Transmission Type (manual/automatic): ");
        String truckTransmission = scanner.nextLine();

        Truck truck = new Truck(truckMake, truckModel, truckYear);
        truck.setCargoCapacity(truckCargoCapacity);
        truck.setTransmissionType(truckTransmission);

        // Display details
        System.out.println("\nVehicle Details:");
        System.out.println("Car:");
        System.out.println("Make: " + car.getMake());
        System.out.println("Model: " + car.getModel());
        System.out.println("Year of Manufacture: " + car.getYearOfManufacture());
        System.out.println("Number of Doors: " + car.getNumberOfDoors());
        System.out.println("Fuel Type: " + car.getFuelType());

        System.out.println("\nMotorcycle:");
        System.out.println("Make: " + motorcycle.getMake());
        System.out.println("Model: " + motorcycle.getModel());
        System.out.println("Year of Manufacture: " + motorcycle.getYearOfManufacture());
        System.out.println("Number of Wheels: " + motorcycle.getNumberOfWheels());
        System.out.println("Motorcycle Type: " + motorcycle.getMotorcycleType());

        System.out.println("\nTruck:");
        System.out.println("Make: " + truck.getMake());
        System.out.println("Model: " + truck.getModel());
        System.out.println("Year of Manufacture: " + truck.getYearOfManufacture());
        System.out.println("Cargo Capacity: " + truck.getCargoCapacity() + " tons");
        System.out.println("Transmission Type: " + truck.getTransmissionType());

        scanner.close();
    }
}