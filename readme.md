
### Summary for Busy Developers

- **Create an Article**: Use `article.NewArticle()` to initialize a new article.
- **Minimal Invariant**: Ensure fields `ID`, `Title`, `Content`, `TextContent`, and `PublishDate` are provided.
- **Normalize and Validate**: Call `article.Normalize()` to trim and validate all fields.
- **Field Limits**: Text fields are trimmed to specific lengths, and URLs are validated with a max length of 4096 characters.
- **Recommended Practice**: Always use the constructor `article.NewArticle()` to ensure the structure is close to its minimal invariant.

### Example Code

```go
package main

import (
    "github.com/editorpost/spider/extract/article"
    "time"
)

func main() {
    // Create a new article
    art := article.NewArticle()

    // Set required fields
    art.Title = "Sample Title"
    art.Content = "This is the content of the article."
    art.TextContent = "This is the text content of the article."
    art.PublishDate = time.Now()

    // Normalize and validate the article
    art.Normalize()
}
```

This documentation provides a comprehensive guide to using the `article` package, covering architecture, usage, and validation limits. By following these guidelines, developers can ensure that their articles are well-structured and validated.
## Article Package Documentation

### Overview

The `article` package provides a structured way to handle and validate article data, ensuring consistency and integrity. The package is designed to normalize and validate various types of content associated with an article, including images, videos, quotes, and social media profiles.

### Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Validation Limits](#validation-limits)
- [Normalization Approach](#normalization-approach)
- [Usage](#usage)
    - [Creating an Article](#creating-an-article)
    - [Minimal Invariant](#minimal-invariant)
    - [Normalization and Validation](#normalization-and-validation)
- [Fields](#fields)
    - [Article](#article)
    - [Image](#image)
    - [Video](#video)
    - [Quote](#quote)
    - [SocialProfile](#socialprofile)
- [Summary for Busy Developers](#summary-for-busy-developers)

### Architecture

The `article` package is built around the `Article` struct, which includes various fields to store article metadata and content. Each nested structure (`Image`, `Video`, `Quote`, and `SocialProfile`) has its own validation and normalization logic to ensure data integrity.

### Validation Limits

The package enforces several validation limits to ensure data consistency and prevent overflow attacks:

- **URL Fields**: Maximum length of 4096 characters.
- **Text Fields**: Trimmed and limited to specific lengths (e.g., title: 255 characters, caption: 500 characters).
- **Author Name**: Limited to 255 characters.
- **Language Code**: Must be a valid ISO 639-1 code (2 characters).
- **Content**: Text content fields are limited to 65000 characters.

These limits ensure that the data remains manageable and secure, suitable for database storage and processing.

### Normalization Approach

Normalization in the `article` package involves:

1. Trimming leading and trailing whitespace from all text fields.
2. Trimming text fields to their maximum allowed lengths.
3. Validating URLs and setting invalid fields to their zero values.
4. Logging validation errors without throwing exceptions, ensuring robustness.

### Usage

#### Creating an Article

To create an article, use the `NewArticle` constructor to initialize a new `Article` struct with default values. This ensures the structure is close to its minimal invariant.

```go
article := article.NewArticle()
```

#### Minimal Invariant

The minimal invariant for an `Article` includes the following required fields:

- `ID`
- `Title`
- `Content`
- `TextContent`
- `PublishDate`

Here is an example of a minimal invariant in JSON format:

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "Sample Title",
  "html": "This is the content of the article.",
  "text": "This is the text content of the article.",
  "published": "2024-01-01T00:00:00Z"
}
```

#### Normalization and Validation

To normalize and validate an article, call the `Normalize` method on the `Article` struct. This method trims and validates all fields, logging any validation errors and clearing invalid fields.

```go
article.Normalize()
```

### Fields

#### Article

The `Article` struct includes the following fields:

- **ID**: UUID of the article (required, max length: 36).
- **Title**: Title of the article (required, max length: 255).
- **Byline**: Author(s) of the article (optional, max length: 255).
- **Content**: Full content of the article (required, max length: 65000).
- **TextContent**: Text content of the article (required, max length: 65000).
- **Excerpt**: Short excerpt of the article (optional, max length: 500).
- **PublishDate**: Publication date of the article (required).
- **ModifiedDate**: Last modification date of the article (optional).
- **Images**: List of images associated with the article.
- **Videos**: List of videos associated with the article.
- **Quotes**: List of quotes associated with the article.
- **Tags**: List of tags associated with the article.
- **Source**: Source URL of the article (optional, max length: 4096).
- **Language**: Language code of the article (required, max length: 2).
- **Category**: Category of the article (optional, max length: 255).
- **SiteName**: Site name where the article is published (optional, max length: 255).
- **AuthorSocialProfiles**: List of social media profiles of the authors.

#### Image

The `Image` struct includes the following fields:

- **URL**: URL of the image (required, max length: 4096).
- **AltText**: Alternative text for the image (optional, max length: 255).
- **Width**: Width of the image in pixels (optional, min: 0).
- **Height**: Height of the image in pixels (optional, min: 0).
- **Caption**: Caption for the image (optional, max length: 500).

#### Video

The `Video` struct includes the following fields:

- **URL**: URL of the video (required, max length: 4096).
- **EmbedCode**: Embed code for the video (optional, max length: 65000).
- **Caption**: Caption for the video (optional, max length: 500).

#### Quote

The `Quote` struct includes the following fields:

- **Text**: Text of the quote (required).
- **Author**: Author of the quote (optional, max length: 255).
- **Source**: Source URL of the quote (optional, max length: 4096).
- **Platform**: Platform where the quote was found (optional, max length: 255).

#### SocialProfile

The `SocialProfile` struct includes the following fields:

- **Platform**: Platform name (required, max length: 255).
- **URL**: URL of the social profile (required, max length: 4096).
