import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HostListComponent } from './host-list.component';

describe('HostListComponent', () => {
  let component: HostListComponent;
  let fixture: ComponentFixture<HostListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HostListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HostListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
