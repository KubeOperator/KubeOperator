import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HostGrantComponent } from './host-grant.component';

describe('HostGrantComponent', () => {
  let component: HostGrantComponent;
  let fixture: ComponentFixture<HostGrantComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HostGrantComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HostGrantComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
