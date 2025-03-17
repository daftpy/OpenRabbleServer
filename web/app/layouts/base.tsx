import { Box, Container, Heading } from "@radix-ui/themes";
import { Link, Outlet } from "react-router";

export default function BaseLayout() {
  return (
    <Container px={"6"}>
      <Box pt={"6"} pb={"4"}>
        <Heading size={"8"} weight={"bold"} className="text-xl pb-1" style={{ color: "var(--slate-12)" }}>
          <Link to="/">OnRabble Server</Link>
        </Heading>
      </Box>
      <Outlet />
    </Container>
  )
}
