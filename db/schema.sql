
CREATE TABLE IF NOT EXISTS event (
    id       bigserial PRIMARY KEY,

    -- The public ID used to reference an event.
    public   bigint UNIQUE NOT NULL,

    -- StoryDevs slug.
    slug     text UNIQUE NOT NULL,

    -- Name and summary to contextualise when listing.
    name     text NOT NULL,
    summary  text NOT NULL,

    start    bigint NOT NULL,
    finish   bigint,

    -- ID of the Discord user who added the event.
    added    text NOT NULL,

    -- Event happens in user's local time.
    local    boolean NOT NULL,

    -- Event recurs each week.
    weekly   boolean NOT NULL,

    -- Categories are defined by StoryDevs.
    category text[] NOT NULL,

    -- Categories are defined by StoryDevs.
    setting text[] NOT NULL
);
