import { CropIcon, Cross1Icon, MagnifyingGlassIcon } from "@radix-ui/react-icons";
import { Badge, Button, DropdownMenu, Flex, Text, TextField } from "@radix-ui/themes";
import { useReducer } from "react";

function reducer(state : any, action: any) {
  switch (action.type) {
    case "set_keyword":
      return {...state, keyword: action.keyword}
    case "add_filter":
      console.log("Adding filter:", action.filter);
      return {...state, activeFilters: [...state.activeFilters, action.filter]}
    case "remove_filter":
      return {...state, activeFilters: [...state.activeFilters.filter((filter : any ) => filter != action.filter)]}
    default:
      return state;
  }
}

export function MessageSearchInput({ 
  keyword,
  filters,
  handleSearch 
} : {
  keyword: string,
  filters: string[],
  handleSearch: ({ filters, keyword }: { filters: any; keyword: string }) => void;
}) {
  const [state, dispatch] = useReducer(reducer, {keyword: keyword, activeFilters: []});

  return (
    <Flex direction={"column"} gap={"1"}>
      <Flex gap={"4"} py={"2"}>
        <TextField.Root
            placeholder="keyword"
            className="grow"
            value={state.keyword}
            onChange={(e) => dispatch({type: "set_keyword", keyword: e.target.value})}
          />
          <DropdownMenu.Root>
            <DropdownMenu.Trigger>
              <Button color="amber">Filter</Button>
            </DropdownMenu.Trigger>
            <DropdownMenu.Content>
              {filters && filters.map((filter: any) => (
                <DropdownMenu.Item key={filter.name} onClick={(e) => dispatch({ type: "add_filter", filter: filter })}>
                  { filter.name }
                </DropdownMenu.Item>
              ))}
            </DropdownMenu.Content>
          </DropdownMenu.Root>
          <Button color="blue" onClick={() => handleSearch({ filters: state.activeFilters, keyword: state.keyword })}><MagnifyingGlassIcon />Search</Button>
      </Flex>
      <Flex gap={"2"} align={"baseline"}>
        <Text weight={"bold"} style={{ color: "var(--indigo-12)" }}>Filters: </Text>
        {state.activeFilters.map((filter: any) => (
          <Badge key={filter.name} size={"2"} color="tomato" onClick={(e) => dispatch({type: "remove_filter", filter: filter})}>
            {filter.name} <Cross1Icon />
          </Badge>
        ))}
      </Flex>
    </Flex>
  )
}