CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS employee (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          username VARCHAR(50) UNIQUE NOT NULL,
                          first_name VARCHAR(50),
                          last_name VARCHAR(50),
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);

CREATE TABLE IF NOT EXISTS organization (
                              id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                              name VARCHAR(100) NOT NULL,
                              description TEXT,
                              type organization_type,
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organization_responsible (
                                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                          organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
                                          user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tenders (
                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                         name VARCHAR(100) NOT NULL,
                         description TEXT,
                         status VARCHAR(20) CHECK (status IN ('Created', 'Published', 'Closed')),
                         version INT DEFAULT 1,
                         organization_id uuid REFERENCES organization(id),
                         creator_id uuid REFERENCES employee(id),
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         service_type VARCHAR(30) CHECK (service_type in ('Construction', 'Delivery', 'Manufacture'))
);

CREATE TABLE IF NOT EXISTS tender_versions (
                                 id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                 tender_id uuid NOT NULL REFERENCES tenders(id) ON DELETE CASCADE,
                                 version INT NOT NULL,
                                 name VARCHAR(100) NOT NULL,
                                 description TEXT,
                                 status VARCHAR(20) CHECK (status IN ('Created', 'Published', 'Closed')),
                                 organization_id uuid NOT NULL,
                                 creator_id uuid NOT NULL,
                                 created_at TIMESTAMP NOT NULL,
                                 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                 service_type text not null CHECK (service_type in ('Construction', 'Delivery', 'Manufacture')),
                                 UNIQUE(tender_id, version)

);

CREATE TABLE IF NOT EXISTS bids (
                      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                      name VARCHAR(100) NOT NULL,
                      description TEXT,
                      status VARCHAR(20) NOT NULL CHECK (status IN ('Created', 'Published', 'Canceled')),
                      version INT NOT NULL DEFAULT 1,
                      tender_id UUID NOT NULL REFERENCES tenders(id) ON DELETE CASCADE ,
                      organization_id UUID NOT NULL REFERENCES organization(id),
                      author_id UUID NOT NULL,
                      author_type VARCHAR(50) NOT NULL CHECK(author_type IN ('Organization', 'User')),
                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS bid_versions (
                              id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                              bid_id uuid not null references bids(id) ON DELETE CASCADE,
                              name VARCHAR(100) NOT NULL,
                              description TEXT,
                              status VARCHAR(20) NOT NULL CHECK (status IN ('Created', 'Published', 'Canceled')),
                              version INT NOT NULL DEFAULT 1,
                              tender_id UUID NOT NULL REFERENCES tenders(id) ON DELETE CASCADE,
                              organization_id UUID NOT NULL REFERENCES organization(id),
                              author_id UUID NOT NULL,
                              author_type VARCHAR(50) NOT NULL CHECK(author_type IN ('Organization', 'User')),
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              UNIQUE(bid_id, version)
);

CREATE TABLE IF NOT EXISTS bid_decisions (
                               id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                               bid_id UUID NOT NULL REFERENCES bids (id) ON DELETE CASCADE,
                               user_id UUID NOT NULL REFERENCES employee (id) ON DELETE CASCADE,
                               decision VARCHAR(20) NOT NULL,
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               UNIQUE (bid_id, user_id)
);

CREATE TABLE IF NOT EXISTS bid_reviews (
                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                             bid_id UUID NOT NULL REFERENCES bids (id) ON DELETE CASCADE,
                             description TEXT NOT NULL CHECK (length(description) <= 1000),
                             author_id UUID NOT NULL REFERENCES employee (id) ON DELETE CASCADE,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);