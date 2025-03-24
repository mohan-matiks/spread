import { z } from "zod";

// Schema for app operating system
export const AppOsSchema = z.enum(["ios", "android"]);

// Schema for creating a new app
export const CreateAppSchema = z.object({
  appName: z
    .string()
    .min(1, "App name is required")
    .max(50, "App name must be less than 50 characters")
    .refine((name) => !/^\s*$/.test(name), "App name cannot be empty"),
  os: AppOsSchema,
});

// Type inference from the schema
export type CreateAppRequest = z.infer<typeof CreateAppSchema>;

// Schema for app response
export const AppSchema = z.object({
  id: z.string(),
  name: z.string(),
  os: z.string(),
  createdAt: z.string().optional(),
  updatedAt: z.string().optional(),
});

export type App = z.infer<typeof AppSchema>;

// Validation function with detailed error handling
export const validateCreateAppRequest = (
  data: unknown
): {
  success: boolean;
  data?: CreateAppRequest;
  error?: Record<string, string>;
} => {
  try {
    const result = CreateAppSchema.parse(data);
    return { success: true, data: result };
  } catch (error) {
    if (error instanceof z.ZodError) {
      // Format errors into a more usable structure
      const formattedErrors: Record<string, string> = {};
      error.errors.forEach((err) => {
        const path = err.path.join(".");
        formattedErrors[path] = err.message;
      });
      return { success: false, error: formattedErrors };
    }
    return { success: false, error: { _form: "Invalid data provided" } };
  }
};
