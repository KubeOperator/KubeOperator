import { ComponentFixture, TestBed } from '@angular/core/testing';

import { IpPoolDeleteComponent } from './ip-pool-delete.component';

describe('IpPoolDeleteComponent', () => {
  let component: IpPoolDeleteComponent;
  let fixture: ComponentFixture<IpPoolDeleteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ IpPoolDeleteComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(IpPoolDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
