import { ArchiveIcon, CrossCircledIcon, MagnifyingGlassIcon, PersonIcon } from "@radix-ui/react-icons";
import { Box, Button, DropdownMenu, Flex, Heading, Text, TextField } from "@radix-ui/themes";
import { Link, useNavigate } from "react-router";
import type { User } from "~/types/components/users";

type props = {
  users: User[];
};

export function UsersPage({ users } : props) {
  const navigate = useNavigate();
  return (
    <>
      <Box py={"1"}>
        <Heading color="indigo">User Management</Heading>
        <Text>Manage your users here. You can ban users, manage their roles, or inspect their activity.</Text>
      </Box>
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
      <Flex pt={"3"} gap={"2"}>
        <Text weight={"bold"} size={"1"}>Searching by: </Text><Text color="amber" size="1">username</Text>
      </Flex>
      <Flex pt={"4"}><Button size={"1"} color="red" onClick={() => navigate("bans")}><ArchiveIcon /> Ban Records</Button></Flex>
      <Flex direction={"column"} pt={"4"}>
        {users.map((user: User, index: number) => (
          <Flex align={"center"} key={index}
            style={{borderBottom: index == users.length - 1 ? "none" : "2px solid var(--indigo-3)"}}
          >
            <Flex direction={"column"} gap={"1"} py={"3"} flexGrow={"1"}>
              <Flex gap={"2"} align={"center"}>
                { user.is_banned ? <CrossCircledIcon color="red" /> : <PersonIcon />}
                <Text weight={"bold"} style={{color: user.is_banned ? "var(--red-9)" : "var(--indigo-12)"}}>
                  <Link to={`/users/profile/${user.username}`}>{user.username}</Link>
                  </Text>
              </Flex>
              <Flex gap={"2"}>
                <Text weight={"bold"} size={"1"}>ID</Text>
                <Text size={"1"} style={{color: "var(--muted-text-color)"}}>{user.id}</Text>
              </Flex>
            </Flex>
            <Button size={"1"} color="iris" onClick={() => navigate(`/users/profile/${user.username}`)}>Manage</Button>
          </Flex>
        ))}
      </Flex>
    </>
  )
}
