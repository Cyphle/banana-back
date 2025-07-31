import { Entity, PrimaryGeneratedColumn, Column } from "typeorm"

@Entity({ name: 'users' })
export class User {

    @PrimaryGeneratedColumn()
    id!: number;

    @Column({ name: 'firstname' })
    firstName!: string;

    @Column({ name: 'lastname' })
    lastName!: string;

    @Column()
    age!: number;

}
