import { FC, useMemo, HTMLAttributes, CSSProperties } from "react";
import { useSpring, animated } from "@react-spring/web";

// https://github.com/JoseRFelix/react-toggle-dark-mode/blob/master/src/index.tsx
// modify to avoid deprecated sub dependencies

export const defaultProperties = {
  dark: {
    circle: {
      r: 9,
    },
    mask: {
      cx: "50%",
      cy: "23%",
    },
    svg: {
      transform: "rotate(40deg)",
    },
    lines: {
      opacity: 0,
    },
  },
  light: {
    circle: {
      r: 5,
    },
    mask: {
      cx: "100%",
      cy: "0%",
    },
    svg: {
      transform: "rotate(90deg)",
    },
    lines: {
      opacity: 1,
    },
  },
  springConfig: { mass: 2.5, tension: 350, friction: 35 },
};

type SVGProps = Omit<HTMLAttributes<HTMLOrSVGElement>, "onChange">;

export interface Props extends SVGProps {
  onChange?: (checked: boolean) => void;
  checked?: boolean;
  style?: CSSProperties;
  size?: number | string;
  animationProperties?: typeof defaultProperties;
  moonColor?: string;
  sunColor?: string;
}

export const DarkModeSwitch: FC<Props> = ({
  onChange,
  children,
  checked = false,
  size = 24,
  animationProperties = defaultProperties,
  moonColor = "white",
  sunColor = "black",
  style,
  ...rest
}) => {
  const uniqueMaskId = useMemo(() => `circle-mask-${Math.random()}`, []);

  const properties = useMemo(() => {
    if (animationProperties !== defaultProperties) {
      return Object.assign(defaultProperties, animationProperties);
    }

    return animationProperties;
  }, [animationProperties]);

  const { circle, svg, lines, mask } = properties[checked ? "dark" : "light"];

  const svgContainerProps = useSpring({
    ...svg,
    config: animationProperties.springConfig,
  });
  const centerCircleProps = useSpring({
    ...circle,
    config: animationProperties.springConfig,
  });
  const maskedCircleProps = useSpring({
    ...mask,
    config: animationProperties.springConfig,
  });
  const linesProps = useSpring({
    ...lines,
    config: animationProperties.springConfig,
  });

  const toggle = () => onChange?.(!checked);

  return (
    // type error for element children, no fix method
    // @ts-ignore https://github.com/pmndrs/react-spring/issues/508
    <animated.svg
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      color={checked ? moonColor : sunColor}
      fill="none"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      stroke="currentColor"
      onClick={toggle}
      style={{
        cursor: "pointer",
        ...svgContainerProps,
        ...style,
      }}
      {...rest}
    >
      <mask id={uniqueMaskId}>
        <rect x="0" y="0" width="100%" height="100%" fill="white" />
        <animated.circle
          // @ts-expect-error
          style={maskedCircleProps}
          r="9"
          fill="black"
        />
      </mask>

      <animated.circle
        cx="12"
        cy="12"
        fill={checked ? moonColor : sunColor}
        // @ts-expect-error
        style={centerCircleProps}
        mask={`url(#${uniqueMaskId})`}
      />

      {/* type error for element children, no fix method */}
      {/* @ts-ignore https://github.com/pmndrs/react-spring/issues/508 */}
      <animated.g stroke="currentColor" style={linesProps}>
        <line x1="12" y1="1" x2="12" y2="3" />
        <line x1="12" y1="21" x2="12" y2="23" />
        <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
        <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
        <line x1="1" y1="12" x2="3" y2="12" />
        <line x1="21" y1="12" x2="23" y2="12" />
        <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
        <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
      </animated.g>
    </animated.svg>
  );
};
export default DarkModeSwitch;
