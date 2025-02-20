import { type RouteConfig, route, index } from "@react-router/dev/routes";

export default [
    index("routes/home.tsx"),
    route("unauthorized", "routes/unauthorized.tsx")
] satisfies RouteConfig;
