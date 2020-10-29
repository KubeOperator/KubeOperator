import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {BaseModelService} from '../../../shared/class/BaseModelService';
import {Observable} from 'rxjs';
import {Notice} from './notice';
import {Page} from '../../../shared/class/Page';

@Injectable({
    providedIn: 'root'
})
export class NoticeService extends BaseModelService<any> {

    baseUrl = '/api/v1/message/mail';
    msgUrl = '/api/v1/message';

    constructor(http: HttpClient) {
        super(http);
    }

    pageBy(page, size): Observable<Page<Notice>> {
        const pageUrl = `${this.baseUrl}?pageNum=${page}&pageSize=${size}`;
        return this.http.get<Page<Notice>>(pageUrl);
    }

    listUnread(): Observable<any> {
        const pageUrl = `${this.msgUrl}/unread`;
        return this.http.get<any>(pageUrl);
    }
}
