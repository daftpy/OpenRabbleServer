
export type AddChannelPayload = {
    name: string;
    description: string | null;
}

export type EditChannelPayload = {
  id: number;
  name: string | null;
  description: string | null;
}

export type ReorderChannelPayload = {
  id: number,
  before_id: number
}

export type ChannelResponse = {
    id: number;
    name: string;
    description: string | null;
};

export type FetchChannelsResponse = {
    channels: ChannelResponse[];
};

export type EditChannelResponse = {
    message: string;
}

export type ReorderChannelResponse = {
    message: string;
}

export type DeleteChannelResponse = {
    message: string;
}

export type AddChannelResponse = {
    message: string;
    name: string;
}
