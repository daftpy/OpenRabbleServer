import { type RouteConfig, route, index, prefix, layout } from "@react-router/dev/routes";

export default [
    index("routes/index.tsx"),
    route("messages", "routes/search_messages.tsx"),
    route("channels", "routes/channel.tsx"),
    ...prefix("users", [
        index("routes/users/index.tsx"),
        route("profile/:userId", "routes/users/profile.tsx"),
        layout("layouts/base.tsx", [
            route("bans", "routes/users/bans.tsx")
        ])
    ]),
    route("unauthorized", "routes/unauthorized.tsx")
] satisfies RouteConfig;
