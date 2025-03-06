import { PersonIcon, TimerIcon } from "@radix-ui/react-icons";
import { Flex, Heading, Text } from "@radix-ui/themes";
import { Link } from "react-router";
import { formatDistance, parseISO } from "date-fns";

export type MessageType = {
  username: string;
  channel: string;
  message: string;
  authored_at: string;
  isLast: boolean;
}

export function Message(props : MessageType) {
  const borderStyle = props.isLast ? "none" : "2px solid var(--indigo-3)"
  const sent = formatDistance(parseISO(props.authored_at), new Date(), { addSuffix: true });
  return (
    <Flex direction="column" gap="1" style={{borderBottom: borderStyle}} py={"2"}> 
      <Flex gap={"1"} align={"center"}>
        <Text size="1" weight={"bold"} color="indigo">{ props.channel } #</Text>
        <Flex gap={"1"} align={"center"}><PersonIcon />
          <Heading size="1">
            <Link to={`/users/profile/${props.username}`}>{ props.username }</Link>
          </Heading>
        </Flex>
      </Flex>
      <Text size="2" wrap={"wrap"} className="py-1">{ props.message }</Text>
      <Flex align={"center"} gap={"1"}><TimerIcon style={{color: "var(--blue-8)", fontWeight: "bold"}} /><Text size={"1"} style={{color: "var(--blue-9)"}}>{sent}</Text></Flex>
    </Flex>
  )
}