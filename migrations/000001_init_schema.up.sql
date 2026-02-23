-- =====================================================
-- 1. ADMINS (Authentication)
-- =====================================================
CREATE TABLE admins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL, -- ('admin', 'coach', 'manager', 'viewer')
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE UNIQUE INDEX idx_unique_admins_username ON admins(username) WHERE deleted_at IS NULL;
COMMENT ON TABLE admins IS 'Table for company XYZ admin authentication';


-- =====================================================
-- 2. TEAMS (Master Teams)
-- =====================================================
CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    logo_url TEXT,
    founded_year INTEGER,
    address TEXT,
    city TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE UNIQUE INDEX idx_unique_teams_name ON teams(name) WHERE deleted_at IS NULL;
CREATE INDEX idx_teams_city ON teams(city) WHERE deleted_at IS NULL;
COMMENT ON TABLE teams IS 'Master data of football teams managed by company XYZ';


-- =====================================================
-- 3. PLAYERS (Team Players)
-- =====================================================
CREATE TABLE players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    name TEXT NOT NULL,
    height_cm NUMERIC(5, 2),
    weight_kg NUMERIC(5, 2),
    position VARCHAR(20) NOT NULL, -- 'forward', 'midfielder', 'defender', 'goalkeeper'
    jersey_number INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE UNIQUE INDEX idx_unique_player_jersey_per_team ON players(team_id, jersey_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_players_team ON players(team_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_players_name ON players(name) WHERE deleted_at IS NULL;
COMMENT ON TABLE players IS 'Player information for each team';


-- =====================================================
-- 4. MATCHES (Match Schedule & Results)
-- =====================================================
CREATE TABLE matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    home_team_id UUID NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    away_team_id UUID NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    match_date DATE NOT NULL,
    match_time TIME NOT NULL,
    home_score INTEGER DEFAULT 0,
    away_score INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'scheduled', -- 'scheduled', 'ongoing', 'finished', 'cancelled'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,
    CONSTRAINT chk_no_self_match CHECK (home_team_id != away_team_id)
);

CREATE INDEX idx_matches_home_team ON matches(home_team_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_away_team ON matches(away_team_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_date ON matches(match_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_status ON matches(status, match_date) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_unique_match_schedule ON matches(home_team_id, away_team_id, match_date) WHERE deleted_at IS NULL;


-- =====================================================
-- 5. MATCH_EVENTS (Match Events)
-- =====================================================
CREATE TABLE match_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE RESTRICT,
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    event_minute INTEGER NOT NULL CHECK (event_minute > 0 AND event_minute <= 120),
    event_type VARCHAR(20) NOT NULL, -- 'goal', 'yellow_card', 'red_card', 'own_goal'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX idx_events_match ON match_events(match_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_events_player ON match_events(player_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_events_type ON match_events(event_type) WHERE deleted_at IS NULL;
