import { type RouteConfig, route, index } from "@react-router/dev/routes";

export default [
    index("routes/index.tsx"),
    route("messages", "routes/search_messages.tsx"),
    route("channels", "routes/channel.tsx"),
    route("users", "routes/users.tsx"),
    route("unauthorized", "routes/unauthorized.tsx")
] satisfies RouteConfig;
