import { type RouteConfig, route, index, prefix, layout } from "@react-router/dev/routes";

export default [
    layout("layouts/base.tsx", [
        index("routes/index.tsx"),
        route("messages", "routes/search_messages.tsx"),
        route("channels", "routes/channel.tsx"),
        ...prefix("users", [
            index("routes/users/index.tsx"),
            route("profile/:userId", "routes/users/profile.tsx"),
            route("bans", "routes/users/bans.tsx"),
        ]),
        route("settings", "routes/settings.tsx"),
        route("unauthorized", "routes/unauthorized.tsx"),
    ]),
] satisfies RouteConfig;