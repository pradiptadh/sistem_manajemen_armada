CREATE TABLE IF NOT EXISTS vehicle_locations (
      vehicle_id TEXT NOT NULL,
      latitude DOUBLE PRECISION NOT NULL,
      longitude DOUBLE PRECISION NOT NULL,
      timestamp TIMESTAMP NOT NULL,
      PRIMARY KEY (vehicle_id, timestamp)
);