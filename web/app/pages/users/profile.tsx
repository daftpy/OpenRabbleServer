import { CircleBackslashIcon, PersonIcon, TimerIcon } from "@radix-ui/react-icons";
import { Box, Button, Container, DropdownMenu, Flex, Heading, Text } from "@radix-ui/themes";
import { Link } from "react-router";

export function UserPage({ username, id } : { username: string, id: string}) {
  return (
    <Container p={"6"}>
      <Heading size={"8"} weight={"bold"} className="text-xl pb-1" style={{ color: "var(--indigo-9)" }}>
        <Link to="/">Your OnRabble Server</Link>
      </Heading>
      <Flex direction="column" pt={"4"} gap={"4"}>
        <Box>
        <Heading style={{color: "var(--subheading-color)"}}>User Information</Heading>
        <Text>Here you can set a users roles, inspect their activity, and ban them from the chatserver if needed.</Text>
        </Box>
        <Flex align={"center"} gap={"4"}>
          <Flex direction={"column"} gap={"2"} flexGrow={"1"}>
            <Flex align={"center"} gap={"3"}>
              <PersonIcon style={{width: "1.5em", height: "1.5em", color: "var(--indigo-12)"}} />
              <Heading style={{color: "var(--highlighted-color)"}}>{username}</Heading>
            </Flex>
            <Flex gap={"2"}>
              <Text weight={"bold"} size={"1"}>ID: </Text>
              <Text size={"1"}>{ id }</Text>
            </Flex>
          </Flex>
          <DropdownMenu.Root>
            <DropdownMenu.Trigger>
              <Button color="red">Ban</Button>
            </DropdownMenu.Trigger>
            <DropdownMenu.Content>

                <DropdownMenu.Item>
                  <CircleBackslashIcon /> Permanent
                </DropdownMenu.Item>
                <DropdownMenu.Item>
                  <TimerIcon /> Temporary
                </DropdownMenu.Item>

            </DropdownMenu.Content>
          </DropdownMenu.Root>
        </Flex>
      </Flex>
    </Container>
  )
}