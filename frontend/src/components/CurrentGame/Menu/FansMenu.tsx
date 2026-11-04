import { Link, Table } from '@heroui/react';
import { FANS } from './fans';

export function FansMenu() {
  return (
    <div
      className={
        'overflow-hidden flex flex-col justify-center mb-3 w-full ' +
        '[&_table]:w-full [&_table]:table-fixed [&_:is(th,td)]:box-border [&_:is(th,td)]:px-1.5 ' +
        '[&_:is(th,td):first-child]:w-12 [&_:is(th,td):first-child]:min-w-12 ' +
        '[&_:is(th,td):last-child]:w-12 [&_:is(th,td):last-child]:min-w-12 ' +
        '[&_:is(th,td):nth-child(2)]:w-[calc((100%-6rem)/3)] [&_:is(th,td):nth-child(2)]:min-w-0'
      }
    >
      <Table className="sm:w-fit self-center">
        <Table.ScrollContainer>
          <Table.Content aria-label="Fans list" selectionMode="single">
            <Table.Header>
              <Table.Column isRowHeader>No.</Table.Column>
              <Table.Column>Fan name</Table.Column>
              <Table.Column>Description</Table.Column>
              <Table.Column>Pts.</Table.Column>
            </Table.Header>
            <Table.Body>
              {FANS.map((fan) => (
                <Table.Row
                  key={fan.no}
                  id={fan.no}
                  className="even:bg-default/50"
                >
                  <Table.Cell>{fan.no}</Table.Cell>
                  <Table.Cell>{fan.name}</Table.Cell>
                  <Table.Cell>{fan.description}</Table.Cell>
                  <Table.Cell>{fan.points}</Table.Cell>
                </Table.Row>
              ))}
            </Table.Body>
          </Table.Content>
        </Table.ScrollContainer>
      </Table>
      <div className="text-end my-3">
        <Link
          href="/fans.jpeg"
          target="_blank"
          rel="noopener noreferrer"
          className="text-sm text-foreground"
        >
          Same, but cooler.jpg
          <Link.Icon />
        </Link>
      </div>
    </div>
  );
}
