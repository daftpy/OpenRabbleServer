import { GearIcon, LockClosedIcon, MagnifyingGlassIcon, PersonIcon } from "@radix-ui/react-icons";
import { Box, Button, Container, Flex, Heading } from "@radix-ui/themes";
import { Link, Outlet, useNavigate } from "react-router";

export default function BaseLayout() {
  const navigate = useNavigate();
  return (
    <Container px={"6"}>
      <Box pt={"6"} pb={"4"}>
        <Heading size={"8"} weight={"bold"} className="text-xl pb-1" style={{ color: "var(--slate-12)" }}>
          <Link to="/">OnRabble Server</Link>
        </Heading>
        <Flex gap={"3"} direction={{initial: "column", sm: "row"}} pt={"4"}>
          <Button onClick={() => navigate("/users")} style={{backgroundColor: "var(--menu-button)"}}><PersonIcon /> User Management</Button>
          <Button onClick={() => navigate("/messages")} style={{backgroundColor: "var(--menu-button)"}}><MagnifyingGlassIcon /> Messages</Button>
          <Button onClick={() => navigate("/about")} style={{backgroundColor: "var(--menu-button)"}}><GearIcon /> Settings</Button>
          <Button onClick={() => navigate("/about")} style={{backgroundColor: "var(--menu-button)"}}><LockClosedIcon /> Keycloak</Button>
        </Flex>
      </Box>
      <Outlet />
    </Container>
  )
}
