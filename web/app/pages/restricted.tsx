import { Box, Container, Flex, Text } from "@radix-ui/themes";

export function Restricted() {
  return (
    <Flex width={"100vw"} height={"100vh"} justify={"center"} align={"center"}>
      <Text align="center" className="text-white" color="crimson" size={"8"}>
        Unauthorized
      </Text>
    </Flex>
  )
}