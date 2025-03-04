import { type RouteConfig, route, index } from "@react-router/dev/routes";

export default [
    index("routes/index.tsx"),
    route("about", "routes/about.tsx"),
    route("channels", "routes/channel.tsx"),
    route("unauthorized", "routes/unauthorized.tsx")
] satisfies RouteConfig;
