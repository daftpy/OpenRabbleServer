import { Button, Dialog, Flex, Text, TextField } from "@radix-ui/themes";
import { useEffect, useState } from "react";
import { useFetcher } from "react-router";
import type { RateLimiter } from "~/pages/settings";

type props = {
  rateLimiter: RateLimiter;
  isOpen: boolean;
  setIsOpen: (open: boolean) => void;
}

export function RateLimitDialog(props: props) {
  const [wordLimit, setWordLimit] = useState<string>(props.rateLimiter.message_limit.toString());
  const [windowSeconds, setWindowSeconds] = useState<string>(props.rateLimiter.window_seconds.toString());

  const rateLimitFetcher = useFetcher();

  useEffect(() => {
    props.setIsOpen(false);
  }, [props.rateLimiter])

  const update = () => {
    rateLimitFetcher.submit(
      {
        wordLimit: wordLimit,
        windowSeconds: windowSeconds
      },
      {
        method: "post",
        action: "/settings",
        encType: "application/x-www-form-urlencoded",
      }
    )
  }
  return (
    <Dialog.Root open={props.isOpen} onOpenChange={(open) => props.setIsOpen(open)}>
      <Dialog.Content>
        <Dialog.Title align={"center"}>Rate Limiter Settings</Dialog.Title>
        <Flex direction={"column"} gap={"3"}>
          <Flex direction={"column"} gap={"1"}>
            <Text>Word Limit</Text>
            <TextField.Root value={wordLimit} onChange={(e) => setWordLimit(e.target.value)} />
          </Flex>
          <Flex direction={"column"} gap={"1"}>
            <Text>WindowSeconds</Text>
            <TextField.Root value={windowSeconds} onChange={(e) => setWindowSeconds(e.target.value)} />
          </Flex>
          <Button onClick={() => update()}>Update</Button>
        </Flex>
      </Dialog.Content>
    </Dialog.Root>
  )
}