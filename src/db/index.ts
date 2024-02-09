import { Database } from 'bun:sqlite';
import { EmployeeCacheI, UserI } from './types';
import { CurrentKeyboardVersion } from '../const';


const db = new Database('db.sqlite', {
    create: true,
});

db.run('CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, messages INTEGER DEFAULT 0, firstname TEXT, lastname TEXT, username TEXT, faculty TEXT, educationForm TEXT, course TEXT, studyGroup TEXT)');
db.run('CREATE TABLE IF NOT EXISTS employeeCache (id INTEGER PRIMARY KEY AUTOINCREMENT, date INTEGER, employees TEXT)');

if (!db.query('PRAGMA table_info(users)').all().map((value: any) => value.name).includes('isAdmin')) {
    db.run('ALTER TABLE users ADD isAdmin BOOLEAN NOT NULL DEFAULT false')
}
if (!db.query('PRAGMA table_info(users)').all().map((value: any) => value.name).includes('keyboardVersion')) {
    db.run('ALTER TABLE users ADD keyboardVersion INTEGER NOT NULL DEFAULT 1')
}

export class User {
    private static readonly createQuery = db.query('INSERT INTO users (id, firstname, lastname, username, keyboardVersion) VALUES($id, $firstname, $lastname, $username, $keyboardVersion) RETURNING *');
    private static readonly findQuery = db.query('SELECT * FROM users WHERE id = $id');
    private static readonly findByUsernameQuery = db.query('SELECT * FROM users WHERE username = $username');
    private static readonly findAllQuery = db.query('SELECT * FROM users');
    private static readonly incrementMessagesQuery = db.query('UPDATE users SET messages = messages + 1 WHERE id = $id');
    private static readonly setInfoQuery = db.query('UPDATE users SET faculty = $faculty, educationForm = $educationForm, course = $course WHERE id = $id');
    private static readonly setStudyGroupQuery = db.query('UPDATE users SET studyGroup = $studyGroup WHERE id = $id');
    private static readonly setIsAdminQuery = db.query('UPDATE users SET isAdmin = $isAdmin WHERE id = $id');
    private static readonly setKeyboardVersionQuery = db.query('UPDATE users SET keyboardVersion = $keyboardVersion WHERE id = $id');

    static create(id: number, firstname: string, lastname?: string, username?: string) {
        const user = this.createQuery.get({
            $id: id,
            $firstname: firstname,
            $lastname: lastname ?? null,
            $username: username ?? null,
            $keyboardVersion: CurrentKeyboardVersion,
        }) as UserI;
        return user;
    }

    static find(id: number) {
        const user = this.findQuery.get({
            $id: id,
        }) as UserI | undefined;
        return user;    
    }

    static findByUsername(username: string) {
        const user = this.findByUsernameQuery.get({
            $username: username,
        }) as UserI | undefined;
        return user;
    }

    static findAll() {
        return this.findAllQuery.all() as UserI[];
    }

    static findOrCreate(id: number, firstname: string, lastname?: string, username?: string) {
        const user = this.find(id);
        return user ?? this.create(id, firstname, lastname, username);
    }

    static incrementMessages(id: number) {
        this.incrementMessagesQuery.run({
            $id: id,
        });
    }

    static setInfo(id: number, faculty?: string, educationForm?: string, course?: string) {
        this.setInfoQuery.run({
            $id: id,
            $faculty: faculty ?? null,
            $educationForm: educationForm ?? null,
            $course: course ?? null,
        });
    }

    static setStudyGroup(id: number, studyGroup?: string) {
        this.setStudyGroupQuery.run({
            $id: id,
            $studyGroup: studyGroup ?? null,
        });
    }

    static reset(id: number) {
        this.setInfo(id);
        this.setStudyGroup(id);
    }

    static setIsAdmin(id: number, value: boolean) {
        this.setIsAdminQuery.run({
            $id: id,
            $isAdmin: value,
        })
    }

    static setKeyboardVersion(id: number, value: number) {
        this.setKeyboardVersionQuery.run({
            $id: id,
            $keyboardVersion: value,
        });
    }

}

export class EmployeeCache {

    private static readonly clearQuery = db.query('DELETE FROM employeeCache');
    private static readonly setQuery = db.query('INSERT INTO employeeCache (date, employees) VALUES($date, $employees)');
    private static readonly getQuery = db.query('SELECT * FROM employeeCache LIMIT 1');


    static set(data: Omit<EmployeeCacheI, 'id'>) {
        this.clearQuery.run();
        this.setQuery.run({
            $date: data.date.getTime(),
            $employees: JSON.stringify(data.employees)
        });
    }

    static get(): EmployeeCacheI | undefined {
        const cache = this.getQuery.get() as { id: number, date: number, employees: string };
        if (!cache) return;
        return {
            id: cache.id,
            date: new Date(cache.date),
            employees: JSON.parse(cache.employees),
        }
    }
}