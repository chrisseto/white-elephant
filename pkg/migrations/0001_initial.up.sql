CREATE TABLE rooms (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  code TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  state TEXT NOT NULL DEFAULT 'WAITING' CHECK (state IN ('WAITING', 'DONE', ''))
);

CREATE TABLE participants (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  name TEXT NOT NULL,
  address TEXT NOT NULL
);
