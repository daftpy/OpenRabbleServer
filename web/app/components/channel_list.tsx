import { Flex, Button, Table, Text } from "@radix-ui/themes";

export interface Channel {
  name: string;
  description: string | null;
}

export default function ChannelList({ channels }  : { channels: Channel[] }) {
  return (
    <Flex flexGrow={"1"}>
      <Table.Root layout={"fixed"}>
        <Table.Header>
          <Table.Row>
            <Table.ColumnHeaderCell width="125px">Channel</Table.ColumnHeaderCell>
            <Table.ColumnHeaderCell width={"auto"}>Description</Table.ColumnHeaderCell>
          </Table.Row>
        </Table.Header>
        <Table.Body>
        { channels && channels.map((channel, index) => (
          <Table.Row key={index}>
            <Table.RowHeaderCell justify="start">{channel.name}</Table.RowHeaderCell>
            <Table.Cell justify={"start"}>
              <Flex align={"center"}>
                <Text className="grow" truncate>{ channel.description ? <>{channel.description}</> : <>...</>}</Text>
                <Button color="red" size={"1"} radius="full" style={{ boxShadow: "var(--shadow-1)" }}>x</Button>
              </Flex>
            </Table.Cell>
          </Table.Row>
        ))}
        </Table.Body>
        </Table.Root>
    </Flex>
  )
}