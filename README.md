# raster-forge

High-performance document normalizer and rasterization engine with multi-resolution pyramid generation for VLM processing.

## Overview

**raster-forge** transforms heterogeneous document inputs (PDFs, images, mixed formats) into unified, rasterized outputs with multi-scale representations optimized for Vision Language Model (VLM) inference.

### Key Features

- **Universal Input Support**: Process single/multi-page PDFs, images (PNG, JPG), and mixed document folders
- **Unified Output**: Generate clean, rasterized PDFs with embedded images (no text layers)
- **Multi-Resolution Pyramid**: Create cascading downsampled versions (300/150/75 DPI) for VLM experimentation
- **WebP Optimization**: Fast encoding with superior compression for pyramid storage
- **Parallel Processing**: Goroutine-based concurrent document and page processing
- **Immutable Inputs**: Source files remain untouched; all outputs in separate directory structure

## Architecture

### Input Structure

```
input/
├── document1.pdf          # Single PDF (single or multi-page)
├── document2.pdf
├── document3/             # Multi-file document
│   ├── page1.pdf
│   ├── page2.pdf
│   └── page3.jpg
└── document4/             # Image-only document
    ├── scan1.png
    └── scan2.png
```

### Output Structure

```
output_YYYYMMDD_HHMMSS/
├── pdfs/
│   ├── document1.pdf      # Rasterized @ 300 DPI
│   ├── document2.pdf
│   ├── document3.pdf
│   └── document4.pdf
└── pyramid/
    ├── document1/
    │   ├── 300dpi/
    │   │   ├── page1.webp
    │   │   └── page2.webp
    │   ├── 150dpi/
    │   │   ├── page1.webp
    │   │   └── page2.webp
    │   └── 75dpi/
    │       ├── page1.webp
    │       └── page2.webp
    └── document2/
        └── ...
```

### Processing Pipeline

1. **Input Scanning**: Identify documents (files or folders)
2. **Normalization**: Merge multi-file documents, handle mixed formats
3. **Rasterization**: Convert all pages to 300 DPI images
4. **PDF Assembly**: Create unified PDF from rasterized pages
5. **Pyramid Generation**: Cascade downsample to 150 DPI → 75 DPI
6. **WebP Encoding**: Save pyramid levels with optimized compression

## Use Cases

### VLM Multi-Scale Inference

Test document recognition accuracy vs. performance across resolutions:
- **300 DPI**: Fine details, small text, complex tables
- **150 DPI**: Balanced quality/speed for general content
- **75 DPI**: Fast layout analysis and document segmentation

### Document Anonymization

Remove text layers and searchable content while preserving visual appearance for privacy-sensitive workflows.

### Format Standardization

Normalize diverse input formats into consistent rasterized PDFs for downstream processing pipelines.

## Technology Stack

- **pdfcpu**: Pure Go PDF manipulation and image extraction
- **golang.org/x/image**: Image processing and downsampling
- **chai2010/webp** or **nativewebp**: WebP encoding with quality control
- **Go 1.21+**: Goroutines for parallel processing

## Roadmap

- [ ] Initial project structure and CLI skeleton
- [ ] Input folder scanning and document detection
- [ ] PDF normalization and merging logic
- [ ] Rasterization engine with configurable DPI
- [ ] Pyramid generation with downsampling algorithms
- [ ] WebP encoding integration
- [ ] Output folder structure creation
- [ ] Parallel processing optimization
- [ ] Configuration file support (YAML/JSON)
- [ ] Docker containerization
- [ ] REST API for integration with orchestrators

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.

## Status

🚧 **In Development** - Initial commit: December 31, 2025

Commercial Go development experience starts: **2025**
