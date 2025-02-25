import { Flex, Button, Table } from "@radix-ui/themes";

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
            <Table.ColumnHeaderCell></Table.ColumnHeaderCell>
          </Table.Row>
        </Table.Header>
        <Table.Body>
        { channels && channels.map((channel, index) => (
          <Table.Row key={index}>
            <Table.RowHeaderCell justify="start">{channel.name}</Table.RowHeaderCell>
            <Table.Cell justify={"start"}>{ channel.description ? <>{channel.description}</> : <>...</>}</Table.Cell>
            <Table.Cell justify="end"><Button color="red" size={"1"} radius="full"  style={{ boxShadow: "var(--shadow-1)" }}>x</Button></Table.Cell>
          </Table.Row>
        ))}
        </Table.Body>
        </Table.Root>
    </Flex>
  )
}