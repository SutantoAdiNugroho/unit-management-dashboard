CREATE TABLE units (
    id VARCHAR(36) PRIMARY KEY,    
    name VARCHAR(255) NOT NULL,
    type ENUM('capsule', 'cabin') NOT NULL,
    status ENUM('Available', 'Occupied', 'Cleaning In Progress', 'Maintenance Needed') NOT NULL,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);