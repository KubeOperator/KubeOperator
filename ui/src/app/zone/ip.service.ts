import {Injectable} from '@angular/core';
import * as ipaddr from 'ipaddr.js';

@Injectable({
  providedIn: 'root'
})
export class IpService {

  constructor() {
  }

  static isBiggerThan(ip1: string, ip2: string): boolean {
    const ip_1 = ipaddr.parse(ip1);
    const ip_2 = ipaddr.parse(ip2);
    return ip1 > ip2;
  }

  static isValid(ip): boolean {
    return ipaddr.isValid(ip);
  }
}
