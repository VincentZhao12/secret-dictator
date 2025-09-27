# Secret Hitler UI Components

A collection of reusable React components built with the propaganda-style theme using Tailwind CSS and the "Black Ops One" font.

## Components

### Button
A styled button component with variants and sizes.

```tsx
import { Button } from './components';

<Button variant="primary" size="md" onClick={handleClick}>
  CREATE GAME
</Button>
<Button variant="secondary" size="lg">
  JOIN GAME
</Button>
```

**Props:**
- `variant`: 'primary' | 'secondary' (default: 'primary')
- `size`: 'sm' | 'md' | 'lg' (default: 'md')
- All standard HTML button attributes

### Heading
Dynamic heading component with semantic levels and style variants.

```tsx
import { Heading } from './components';

<Heading level={1} variant="title">SECRET</Heading>
<Heading level={2} variant="subtitle">HITLER</Heading>
<Heading level={3} variant="section">Game Rules</Heading>
```

**Props:**
- `level`: 1 | 2 | 3 | 4 | 5 | 6 (default: 1)
- `variant`: 'title' | 'subtitle' | 'section' (default: 'title')

### Text
Typography component for body text with different variants.

```tsx
import { Text } from './components';

<Text variant="body">Regular text content</Text>
<Text variant="description">Italicized description text</Text>
<Text variant="caption">Small caption text</Text>
<Text variant="footer">Footer text with tracking</Text>
```

**Props:**
- `variant`: 'body' | 'description' | 'caption' | 'footer' (default: 'body')
- `as`: 'p' | 'span' | 'div' (default: 'p')

### Card
Container component with the signature border and shadow styling.

```tsx
import { Card } from './components';

<Card size="md">
  <h1>Card Content</h1>
</Card>
```

**Props:**
- `size`: 'sm' | 'md' | 'lg' (default: 'md')

### Container
Page-level wrapper with the gradient background and layout.

```tsx
import { Container } from './components';

<Container>
  <Card>Page content here</Card>
</Container>
```

### Badge
Small status indicator with color variants.

```tsx
import { Badge } from './components';

<Badge variant="primary">5/10 Players</Badge>
<Badge variant="warning">Waiting</Badge>
<Badge variant="danger">Game Over</Badge>
```

**Props:**
- `variant`: 'primary' | 'secondary' | 'warning' | 'danger' (default: 'primary')
- `size`: 'sm' | 'md' (default: 'md')

### Divider
Decorative line separator.

```tsx
import { Divider } from './components';

<Divider className="w-32 mx-auto" />
<Divider orientation="vertical" className="h-8" />
```

**Props:**
- `orientation`: 'horizontal' | 'vertical' (default: 'horizontal')

## Design System

### Colors
- Primary Orange: `bg-orange-600`, `bg-orange-700`
- Secondary Cream: `bg-cream` (custom CSS variable)
- Accent: `bg-orange-200/90`
- Text: `text-black`, `text-orange-700`
- Borders: `border-black`

### Typography
- Font Family: "Black Ops One" (Google Fonts)
- Base Class: `font-propaganda`
- Weights: `font-extrabold`
- Tracking: `tracking-wider`, `tracking-widest`

### Shadows & Effects
- Drop Shadow: `drop-shadow-[2px_2px_0px_orange]`
- Box Shadow: `shadow-[6px_6px_0px_black]`
- Button Shadow: `shadow-[4px_4px_0px_black]`
- Hover Effects: `hover:scale-[1.03]`

### Layout
- Border Width: `border-4`
- Border Radius: `rounded-xl`, `rounded-sm`
- Padding: Responsive with `p-6` to `p-16`

## Usage

Import components from the central index:

```tsx
import { Button, Heading, Text, Card, Container } from './components';
```

All components are built with TypeScript for type safety and use Tailwind CSS for styling. The design maintains the propaganda poster aesthetic with bold typography, high contrast, and dramatic shadows.