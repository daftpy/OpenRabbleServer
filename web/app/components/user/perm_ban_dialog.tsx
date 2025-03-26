import { Button, Dialog, Flex, Text, TextField } from "@radix-ui/themes";
import { useState } from "react";
import { useFetcher } from "react-router";
import { BanDialog } from "~/pages/users/profile";

type props = {
  id: string;
  username: string;
  dialog?: BanDialog;
  setDialog: (dialog : BanDialog | null) => void;
}

export function PermBanDialog(props: props) {
  const [reason, setReason] = useState<string>("");
  const userFetcher = useFetcher();

  const handleBanUser = () => {
    userFetcher.submit(
      {
        reason: reason,
        banishedId: props.id
      },
      {
        method: "POST",
        action: `/users/profile/${props.username}`,
      }
    )
  }
  return (
    <Dialog.Root open={props.dialog == BanDialog.PermaBan} onOpenChange={() => props.setDialog(null)}>
      <Dialog.Content>
        <Flex direction={"column"} py={"6"} gap={"6"}>
          <Flex direction={"column"} gap={"2"}>
            <Text align={"center"} weight={"bold"}>Permananet Ban</Text>
            <Text size={"2"}>Reason</Text>
            <TextField.Root placeholder="Reason" onChange={(e) => setReason(e.target.value)} />
          </Flex>
          <Button color="red" onClick={() => {handleBanUser()}}>Ban</Button>
        </Flex>
      </Dialog.Content>
    </Dialog.Root>
  )
}