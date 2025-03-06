import { MagnifyingGlassIcon, PersonIcon } from "@radix-ui/react-icons";
import { Box, Button, Container, DropdownMenu, Flex, Heading, Text, TextField } from "@radix-ui/themes";
import { Link, useNavigate } from "react-router";

export function UsersPage({ users } : any) {
  const navigate = useNavigate();
  console.log("USERS", users);
  users.map((user: any) => {
    console.log("A USER:", user);
  });
  return (
    <Container p={"6"}>
      <Heading size={"8"} weight={"bold"} className="text-xl pb-1" style={{ color: "var(--indigo-9)" }}>
        <Link to="/">Your OnRabble Server</Link>
      </Heading>
      <Box pt={"4"}>
      <Heading style={{color: "var(--subheading-color)"}}>User Management</Heading>
      <Text>Manage your users here. You can ban users, manage their roles, or inspect their activity.</Text>
      <Flex pt={"2"} gap={"4"}>
        <TextField.Root placeholder="username" className="grow" />
        <DropdownMenu.Root>
          <DropdownMenu.Trigger>
            <Button color="amber">Search by</Button>
          </DropdownMenu.Trigger>
          <DropdownMenu.Content>
            <DropdownMenu.Item>
              Username
            </DropdownMenu.Item>
            <DropdownMenu.Item>
                ID
              </DropdownMenu.Item>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
        <Button color={"blue"}><MagnifyingGlassIcon />Search</Button>
      </Flex>
      <Flex pt={"2"} gap={"2"}>
        <Text weight={"bold"} size={"1"}>Searching by: </Text><Text color="amber" size="1">username</Text>
      </Flex>
      </Box>
      <Flex direction={"column"} pt={"4"}>
        {users.map((user: any, index: number) => (
          <Flex align={"center"}
            style={{borderBottom: index == users.length - 1 ? "none" : "2px solid var(--indigo-3)"}}
          >
            <Flex direction={"column"} gap={"1"} py={"3"} flexGrow={"1"}>
              <Flex gap={"2"} align={"center"}>
                <PersonIcon />
                <Text weight={"bold"} style={{color: "var(--highlighted-color)"}}>
                  <Link to={`/users/profile/${user.username}`}>{user.username}</Link>
                  </Text>
              </Flex>
              <Flex gap={"2"}>
                <Text weight={"bold"} size={"1"}>ID</Text>
                <Text size={"1"} >{user.id}</Text>
              </Flex>
            </Flex>
            <Button size={"1"} color="iris" onClick={() => navigate(`/users/profile/${user.username}`)}>Manage</Button>
          </Flex>
        ))}
      </Flex>
    </Container>
  )
}