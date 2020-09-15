import { Body, Controller, Delete, Get, Param, Post, Put } from '@nestjs/common';
import { UsersService } from './users.service';
import { CreateUserDto, UpdateUserDto } from 'src/models/dto/users';

@Controller('/users')
export class UsersController {
  constructor(private readonly usersService: UsersService) {
  }

  @Post()
  create(@Body('user') user: CreateUserDto) {
    return this.usersService.createUser(user);
  }

  @Get()
  getUsers() {
    return this.usersService.getAllUsers();
  }

  @Get(':id')
  getUser(@Param('id') id: string) {
    return this.usersService.getUserById(id);
  }

  @Put(':id')
  updateUser(@Param('id') id: string, @Body() user: UpdateUserDto) {
    this.usersService.updateUser(id, user);
    return null;
  }

  @Delete(':id')
  deleteUser(@Param('id') id: string) {
    this.usersService.deleteUser(id);
    return null;
  }
}
