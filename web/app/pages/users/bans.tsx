import { PersonIcon } from "@radix-ui/react-icons";
import { Button, Flex, Heading, Table, Text } from "@radix-ui/themes";
import { formatDistance, parseISO } from "date-fns";
import type { BanRecord } from "~/types/components/users";

type props = {
  records: BanRecord[]
}

export function BansPage(props : props) {
  return (
    <Flex direction={"column"}>
      <Heading color="indigo">Bans</Heading>
      <Text>Here you can review bans and grant pardons.</Text>
      <Table.Root>
        <Table.Header>
          <Table.Row>
            <Table.ColumnHeaderCell>User</Table.ColumnHeaderCell>
            <Table.ColumnHeaderCell>Reason</Table.ColumnHeaderCell>
            <Table.ColumnHeaderCell style={{width: "80px", textWrap: "nowrap"}}>Start</Table.ColumnHeaderCell>
            <Table.ColumnHeaderCell>Duration</Table.ColumnHeaderCell>
            <Table.ColumnHeaderCell>Remaining</Table.ColumnHeaderCell>
            <Table.ColumnHeaderCell style={{width: "80px", textWrap: "nowrap"}} align="right">Pardon</Table.ColumnHeaderCell>
          </Table.Row>
        </Table.Header>

        <Table.Body>
          {props.records.length > 0 && props.records.map((record : BanRecord) => (
            <Table.Row align={"center"} key={record.id}>
              <Table.RowHeaderCell maxWidth={"120px"}>
                <Flex gap={"2"} align={"center"} justify={"start"}>
                  <PersonIcon style={{flexShrink: "0"}} />
                  <Text wrap={"nowrap"} weight={"bold"} truncate>
                    {record.banished_username}
                  </Text>
                </Flex>
              </Table.RowHeaderCell>
              <Table.Cell style={{maxWidth: "200px", overflow: "hidden", textOverflow: "ellipsis"}}>
                <Text wrap={"nowrap"} style={{textOverflow: "ellipsis"}}>
                  This is a really long reason hahahahaha This is a really long reason hahahahaha
                </Text>
              </Table.Cell>
              <Table.Cell  style={{width: "50px", textWrap: "nowrap"}}>
                <Text>
                  {new Date(record.start).toLocaleDateString()}
                </Text>
              </Table.Cell>
              <Table.Cell>
                <Text>{record.end != null ? <>{formatDistance(parseISO(record.end), record.start).replace("about", "")}</> : <>Forever</>}</Text>
              </Table.Cell>
              <Table.Cell>
                {record.end == null ? (
                  // If there is no 'end' property, the ban is permanent
                  <>Forever</>
                ) : new Date(record.end) > new Date() ? (

                  // If the 'end' date has not passed, format how long until the ban expires
                  <Text wrap={"nowrap"}>{formatDistance(parseISO(record.end), new Date())}</Text>
                ) : (

                  // If the 'end' date has passed, the ban has expired
                  <>Expired</>
                )}
              </Table.Cell>
              <Table.Cell style={{width: "30px"}} align="right">
                <Button color="red" size={"1"} disabled={record.end !== null && new Date(record.end) <= new Date()}>
                  Pardon
                </Button>
              </Table.Cell>
            </Table.Row>
          ))}
        </Table.Body>
      </Table.Root>
    </Flex>
  )
}
