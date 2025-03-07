import { Link1Icon, PersonIcon, TimerIcon } from "@radix-ui/react-icons";
import { Box, Flex, Heading, Text } from "@radix-ui/themes";
import { Link } from "react-router";
import { formatDistance, parseISO } from "date-fns";

export type MessageType = {
  username: string;
  channel: string;
  message: string;
  authored_at: string;
}

type props = {
  meessage: MessageType;
  isLast: boolean;
  hidePermaLink: boolean;
}

export function Message(props : props) {
  const borderStyle = props.isLast ? "none" : "2px solid var(--indigo-3)"
  const sent = formatDistance(parseISO(props.meessage.authored_at), new Date(), { addSuffix: true });
  return (
    <Flex direction="column" gap="1" style={{borderBottom: borderStyle}} py={"2"}> 
      <Flex direction={"column"} px={"1"}>
        <Flex gap={"1"} align={"center"}>
          <Text size="1" weight={"bold"} color="indigo">{ props.meessage.channel } #</Text>
          <Flex align={"center"} justify={"between"} flexGrow={"1"}>
            <Flex gap={"1"} align={"center"}><PersonIcon />
              <Heading size="1" style={{color: "var(--link-color)"}}>
                <Link to={`/users/profile/${props.meessage.username}`}>{ props.meessage.username }</Link>
              </Heading>
            </Flex>
            {!props.hidePermaLink && (
              <Box>
                <Link to={`/messages`}><Link1Icon style={{color: "var(--gray-10)", width: "12px", height: "12px", cursor: "pointer"}} /></Link>
              </Box>
            )}
          </Flex>
        </Flex>
        <Text size="2" wrap={"wrap"} className="py-1">{ props.meessage.message }</Text>
        <Flex align={"center"} gap={"1"}><TimerIcon style={{color: "var(--muted-text-color)", fontWeight: "bold", width: "13px", height: "13px"}} /><Text size={"1"} style={{color: "var(--muted-text-color)"}}>{sent}</Text></Flex>
      </Flex>
    </Flex>
  )
}