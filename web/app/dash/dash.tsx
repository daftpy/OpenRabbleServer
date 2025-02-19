import { Button, Flex, Heading, Text, TextField } from "@radix-ui/themes"

export function Dash() {
  return (
    <main className="p-4">
      <Flex direction="column" gap="3">
        <div>
          <Heading className="font-bold text-xl">Your OnRabble Server</Heading>
        </div>
        <Text>Welcome to your dashboard.</Text>
        <div>
          <Heading className="font-bold" color="indigo" style={{ color: "var(--indigo-9)"}}>Channels</Heading>
          <Text m="0">You can add a new channel or manage your channels below.</Text>
        </div>
        <Flex direction="row" gap="4">
          <TextField.Root placeholder="Add a new channel">
          </TextField.Root>
          <Button>Add</Button>
        </Flex>
      </Flex>
    </main>
  )
}