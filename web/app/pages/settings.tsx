import { Link2Icon, StopwatchIcon } from "@radix-ui/react-icons";
import { Box, Button, Flex, Heading, Text} from "@radix-ui/themes";
import { useState } from "react";
import { RateLimitDialog } from "~/components/settings/rate_limit_dialog";

export type RateLimiter = {
  id: number;
  message_limit: number;
  window_seconds: number;
}

type props = {
  rate_limiter: RateLimiter;
}

export function SettingsPage(props: props) {
  const [isOpen, setIsOpen] = useState<boolean>(false);
  return (
    <>
      <Box pt={"2"} pb={"4"}>
        <Heading color="indigo">Settings</Heading>
        <Text>Here you can update server settings such as rate limits, enable/disable registration, or mute chat.</Text>
      </Box>
      <Flex direction={"column"}>
        <Heading size="5" style={{color: "var(--subheading-color)"}}>Rate Limiter</Heading>
        <Text>This restricts how frequently chat messages will be processed by your server for all users.</Text>
        <Box pt={"2"}>
          <Flex gap={"4"} pb={"2"}>
            <Flex gap={"1"} align={"center"}>
              <Link2Icon />
              <Text weight={"bold"} style={{color: "var(--indigo-12)"}}>Message Limit:</Text>
              <Text>{ props.rate_limiter.message_limit }</Text>
            </Flex>
            <Flex gap={"1"} align={"center"}>
              <StopwatchIcon />
              <Text weight={"bold"} style={{color: "var(--indigo-12)"}}>Window Seconds:</Text>
              <Text>{ props.rate_limiter.window_seconds }</Text>
            </Flex>
            <RateLimitDialog isOpen={isOpen} setIsOpen={setIsOpen} rateLimiter={props.rate_limiter} />
          </Flex>
          <Button size={"1"} style={{float: "right"}} onClick={() => setIsOpen(true)}>Edit</Button>
        </Box>
      </Flex>
    </>
  )
}