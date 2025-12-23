import z from "zod";

export const ZTodoComment = z.object({
  id: z.uuid(),
  todoId: z.uuid(),
  userId: z.string(),
  content: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
});