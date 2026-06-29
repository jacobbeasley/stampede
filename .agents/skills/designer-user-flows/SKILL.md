# Designer User Flows Skill

## Overview
This skill provides guidance for creating user flow diagrams and sitemaps using Mermaid notation.

## When to Use
- After design brief is complete
- Before mockup generation
- When documenting page hierarchies and user journeys

## Process

### 1. Sitemap Creation
- Identify all pages/routes from design brief
- Determine page hierarchy and relationships
- Create Mermaid graph diagram
- Document each page's purpose and components

### 2. User Flow Diagrams
- Identify key user journeys (e.g., create, edit, delete flows)
- Model decision points and conditional paths
- Create flowchart diagrams for each key flow
- Document states and transitions

### 3. Navigation Structure
- Define global navigation (appears on all pages)
- Define contextual navigation (page-specific)
- Document breadcrumb logic if applicable

## Output
- Mermaid sitemap diagram (`.agents/plans/designer-output/sitemaps/[project]-sitemap.md`)
- Mermaid user flow diagrams (`.agents/plans/designer-output/sitemaps/[project]-flows.md`)

## Best Practices
- Keep sitemap high-level; don't drill into components
- Use meaningful node labels (not just URLs)
- Group related pages with subgraphs
- Document flows that match primary user goals
- Consider mobile navigation differences early
