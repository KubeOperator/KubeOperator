import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ServiceRouteComponent } from './service-route.component';

describe('ServiceRouteComponent', () => {
  let component: ServiceRouteComponent;
  let fixture: ComponentFixture<ServiceRouteComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ServiceRouteComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ServiceRouteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
