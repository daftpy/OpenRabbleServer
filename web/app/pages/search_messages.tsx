import { Button, DropdownMenu, Flex, Heading, Text, TextField } from "@radix-ui/themes";

export function SearchMessagesPage() {
    return (
        <Flex direction={"column"} maxWidth={"990px"} m={"auto"} px={"4"} py={"6"}>
            <Heading size={"8"} weight={"bold"} className="text-xl pb-1" style={{color: "var(--indigo-9)"}}>Your OnRabble Server</Heading>
            <Text>Search messages stored in your cache or database. Filter by <strong style={{color: "var(--link-color)"}}>channel </strong>or <strong style={{color: "var(--link-color)"}}>keyword</strong>.</Text>
            <Flex py={"4"} gap={"4"}>
            <TextField.Root 
                placeholder="keyword"
                className="grow"
            />
            <DropdownMenu.Root>
                <DropdownMenu.Trigger>
                    <Button>
                        Channel
                    </Button>
                </DropdownMenu.Trigger>
                <DropdownMenu.Content>
                    <DropdownMenu.Item>General</DropdownMenu.Item>
                    <DropdownMenu.Item>Tech</DropdownMenu.Item>
                    <DropdownMenu.Item>News</DropdownMenu.Item>
                </DropdownMenu.Content>
            </DropdownMenu.Root>
            <Button color="grass">Search</Button>
            </Flex>
        </Flex>
    )
}