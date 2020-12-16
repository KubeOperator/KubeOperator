import { ComponentFixture, TestBed } from '@angular/core/testing';

import { IpListComponent } from './ip-list.component';

describe('IpListComponent', () => {
  let component: IpListComponent;
  let fixture: ComponentFixture<IpListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ IpListComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(IpListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
