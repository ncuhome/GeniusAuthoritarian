import { FC, useState } from "react";
import {
  Box,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  SelectChangeEvent,
  Stack,
  Typography,
} from "@mui/material";
import { LineChart } from "@mui/x-charts/LineChart";
import Block from "@components/user/Block";
import useLoginData from "@hooks/data/useLoginData";

const pad = (n: number) => n.toString().padStart(2, "0");

const countByPeriod = (period: Admin.DataFetchRange): [string[], number[]] => {
  const { data: loginData } = useLoginData(period);
  if (!loginData) return [[], []];

  const currentTime = Date.now();
  const periodInMilliseconds = {
    week: 7 * 24 * 60 * 60 * 1000,
    month: 4 * 7 * 24 * 60 * 60 * 1000,
    year: 12 * 4 * 7 * 24 * 60 * 60 * 1000,
  }[period];
  const startTime = currentTime - periodInMilliseconds;
  const endTime = currentTime;

  const timestamps: string[] = [];
  const counts: number[] = [];

  const initData = {
    week: () => {
      for (let t = startTime; t <= endTime; t += 24 * 60 * 60 * 1000) {
        const date = new Date(t);
        timestamps.push(
          `${pad(date.getMonth() + 1)}月${pad(date.getDate())}日`,
        );
        counts.push(0);
      }
    },
    month: () => {
      let currentDate = new Date(startTime);
      while (currentDate <= new Date(endTime)) {
        const nextDate = new Date(
          currentDate.getFullYear(),
          currentDate.getMonth(),
          currentDate.getDate() + 7,
        );
        timestamps.push(
          `${pad(currentDate.getMonth() + 1)}月${pad(currentDate.getDate())}日-${pad(nextDate.getMonth() + 1)}月${pad(nextDate.getDate() - 1)}日`,
        );
        counts.push(0);
        currentDate = nextDate;
      }
    },
    year: () => {
      for (let m = 0; m < 12; m++) {
        timestamps.push(`${pad(m + 1)}月`);
        counts.push(0);
      }
    },
  }[period];
  initData();

  for (const record of loginData.records) {
    const createdAt = record.createdAt * 1000;
    if (createdAt >= startTime && createdAt <= endTime) {
      const date = new Date(createdAt);
      if (period === "week") {
        const index = Math.floor(
          (createdAt - startTime) / (24 * 60 * 60 * 1000),
        );
        counts[index]++;
      } else if (period === "month") {
        const weekIndex = Math.floor((date.getDate() - 1) / 7);
        counts[weekIndex]++;
      } else {
        const monthIndex = date.getMonth();
        counts[monthIndex]++;
      }
    }
  }

  return [timestamps, counts];
};

const Status: FC = () => {
  const [period, setPeriod] = useState<"week" | "month" | "year">("week");
  const [timestamps, counts] = countByPeriod(period);

  const handlePeriodChange = (event: SelectChangeEvent<string>) => {
    setPeriod(event.target.value as "week" | "month" | "year");
  };

  const data = countByPeriod(period);

  return (
    <Block>
      <Stack spacing={2}>
        <Stack flexDirection={"row"}>
          <LineChart
            xAxis={[
              {
                data: data[0],
                label: "日期",
                scaleType: "point",
              },
            ]}
            series={[{ data: data[1], label: "登录次数" }]}
            width={700}
            height={300}
          />
          <Box
            sx={{
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              mt: 0,
            }}
          >
            <FormControl sx={{ m: 1, minWidth: 120 }}>
              <InputLabel id="period-select-label">Period</InputLabel>
              <Select
                labelId="period-select-label"
                value={period}
                onChange={handlePeriodChange}
                label="Period"
              >
                <MenuItem value="week">Week</MenuItem>
                <MenuItem value="month">Month</MenuItem>
                <MenuItem value="year">Year</MenuItem>
              </Select>
            </FormControl>
            <Typography variant="h5" gutterBottom>
              {period.toUpperCase()} counts:
            </Typography>
            <Box>
              {timestamps.map((timestamp, index) => (
                <Typography
                  key={index}
                >{`${timestamp}: ${counts[index]} 次`}</Typography>
              ))}
            </Box>
          </Box>
        </Stack>
      </Stack>
    </Block>
  );
};
export default Status;
