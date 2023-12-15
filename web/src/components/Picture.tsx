import { FC, CSSProperties } from "react";

interface Props {
  name: string;
  alt: string;
  dir?: string;
  defaultType?: string;
  sources?: Source[];

  imgStyle?: CSSProperties;
}

interface Source {
  fileType: string;
  mimeType: string;
}

const Picture: FC<Props> = ({
  name,
  alt,
  dir,
  defaultType = "png",
  sources = [{ fileType: "webp", mimeType: "image/webp" }],
  imgStyle,
}) => {
  const getImageUrl = (name: string, type: string) => {
    if (dir) {
      return new URL(`/src/assets/img/${dir}/${name}.${type}`, import.meta.url)
        .href;
    }
    return new URL(`/src/assets/img/${name}.${type}`, import.meta.url).href;
  };

  return (
    <picture>
      {sources.map((source) => (
        <source
          key={JSON.stringify(source.fileType)}
          srcSet={getImageUrl(name, source.fileType)}
          type={source.mimeType}
        />
      ))}
      <img alt={alt} src={getImageUrl(name, defaultType)} style={imgStyle} />
    </picture>
  );
};
export default Picture;
