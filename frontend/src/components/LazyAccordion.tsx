import { Accordion } from '@heroui/react';
import { ReactNode, useState } from 'react';

type Key = string | number;

export type LazyAccordionSection = {
  id: string;
  title: ReactNode;
  icon?: ReactNode;
  // content renders only while the section is expanded.
  content: () => ReactNode;
};

// LazyAccordion renders a section's content only while that section is open.
// HeroUI v3 Accordion keeps every panel mounted, which would fetch game logs
// and render the long fans table as soon as a menu opens, even when collapsed.
export default function LazyAccordion({
  sections,
  defaultExpandedKeys = [],
  allowsMultipleExpanded = false,
  className,
}: {
  sections: LazyAccordionSection[];
  defaultExpandedKeys?: string[];
  allowsMultipleExpanded?: boolean;
  className?: string;
}) {
  const [expandedKeys, setExpandedKeys] = useState<Set<Key>>(
    () => new Set(defaultExpandedKeys)
  );

  return (
    <Accordion
      className={className}
      allowsMultipleExpanded={allowsMultipleExpanded}
      expandedKeys={expandedKeys}
      onExpandedChange={setExpandedKeys}
    >
      {sections.map((section) => (
        <Accordion.Item key={section.id} id={section.id}>
          <Accordion.Heading>
            <Accordion.Trigger>
              {section.icon ? (
                <span className="me-3 shrink-0">{section.icon}</span>
              ) : null}
              {section.title}
              <Accordion.Indicator />
            </Accordion.Trigger>
          </Accordion.Heading>
          <Accordion.Panel>
            <Accordion.Body>
              {expandedKeys.has(section.id) ? section.content() : null}
            </Accordion.Body>
          </Accordion.Panel>
        </Accordion.Item>
      ))}
    </Accordion>
  );
}
