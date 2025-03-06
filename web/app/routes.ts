import { type RouteConfig, route, index, prefix } from "@react-router/dev/routes";

export default [
    index("routes/index.tsx"),
    route("messages", "routes/search_messages.tsx"),
    route("channels", "routes/channel.tsx"),
    ...prefix("users", [
        index("routes/users.tsx"),
        route(":userId", "routes/user.tsx")
    ]),
    route("unauthorized", "routes/unauthorized.tsx")
] satisfies RouteConfig;
