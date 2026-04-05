CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'user')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Rooms
CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    capacity INTEGER CHECK (capacity > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Schedules
CREATE TABLE IF NOT EXISTS schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id UUID NOT NULL UNIQUE REFERENCES rooms(id) ON DELETE CASCADE,
    days_of_week INTEGER[] NOT NULL,
    start_time TIME WITHOUT TIME ZONE NOT NULL,
    end_time TIME WITHOUT TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT valid_days_of_week CHECK (
        array_length(days_of_week, 1) BETWEEN 1 AND 7
    ),
    CONSTRAINT valid_time_range CHECK (start_time < end_time)
);

-- Slots
CREATE TABLE IF NOT EXISTS slots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    is_available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT slot_duration_check CHECK (
        EXTRACT(EPOCH FROM (end_time - start_time)) = 1800
    ),
    CONSTRAINT valid_time_range CHECK (start_time < end_time)
);

-- Bookings
CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    slot_id UUID NOT NULL REFERENCES slots(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'cancelled')),
    conference_link TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_rooms_name ON rooms(name);
CREATE INDEX idx_schedules_room_id ON schedules(room_id);
CREATE INDEX idx_slots_room_id ON slots(room_id);
CREATE INDEX idx_slots_start_time ON slots(start_time);
CREATE INDEX idx_slots_is_available ON slots(is_available);
CREATE INDEX idx_slots_room_time ON slots(room_id, start_time);
CREATE INDEX idx_bookings_slot_id ON bookings(slot_id);
CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_user_status ON bookings(user_id, status);

-- Partial unique index: prevents multiple active bookings on the same slot
CREATE UNIQUE INDEX unique_active_booking_per_slot ON bookings (slot_id) WHERE (status = 'active');